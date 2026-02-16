from setuptools import setup, find_packages

setup(
    name="sijil",
    version="1.0.0",
    packages=find_packages(),
    install_requires=["requests"],
    description="The Paranoid Log Engine SDK for Python",
    author="Sijil",
    url="https://github.com/sijil-oss/sijil-python",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
    ],
)